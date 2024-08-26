package main

import (
    "context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ec2"
    ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)



func resourceEC2Instance() *schema.Resource {
    return &schema.Resource{
        Create: resourceEC2InstanceCreate,
        Read:   resourceEC2InstanceRead,
        Update: resourceEC2InstanceUpdate,
        Delete: resourceEC2InstanceDelete,

        Schema: map[string]*schema.Schema{
            "instance_type": {
                Type:     schema.TypeString,
                Required: true,
            },
            "ami": {
                Type:     schema.TypeString,
                Required: true,
            },
            "instance_id": {
                Type:     schema.TypeString,
                Computed: true,
            },
        },
    }
}

func resourceEC2InstanceCreate(d *schema.ResourceData, m interface{}) error {
    // Load the default AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
    if err != nil {
        return err
    }

    // Create an EC2 client
    svc := ec2.NewFromConfig(cfg)

    // Convert the instance type from string to ec2types.InstanceType
    instanceType := ec2types.InstanceType(d.Get("instance_type").(string))

    // Prepare the input for the RunInstances API
    runResult, err := svc.RunInstances(context.TODO(), &ec2.RunInstancesInput{
        ImageId:      aws.String(d.Get("ami").(string)),
        InstanceType: instanceType,  // Use the ec2types.InstanceType value
        MinCount:     aws.Int32(1),
        MaxCount:     aws.Int32(1),
    })
    if err != nil {
        return err
    }

    // Retrieve the instance ID and store it in the state
    instanceID := runResult.Instances[0].InstanceId
    d.SetId(*instanceID)
    d.Set("instance_id", instanceID)

    return resourceEC2InstanceRead(d, m)
}

func resourceEC2InstanceRead(d *schema.ResourceData, m interface{}) error {
    // Load the default AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
    if err != nil {
        return err
    }

    // Create an EC2 client
    svc := ec2.NewFromConfig(cfg)

    // Use the DescribeInstances API to retrieve information about the instance
    instanceID := d.Id()
    result, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
        InstanceIds: []string{instanceID},
    })
    if err != nil {
        return err // Return the error directly if there's a problem with the API call
    }

    // Check if the instance was found
    if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
        // Instance not found, possibly terminated; remove it from state
        d.SetId("")
        return nil
    }

    // Update the instance's attributes in the Terraform state
    d.Set("instance_type", result.Reservations[0].Instances[0].InstanceType)
    d.Set("ami", result.Reservations[0].Instances[0].ImageId)

    return nil
}



func resourceEC2InstanceUpdate(d *schema.ResourceData, m interface{}) error {
    return resourceEC2InstanceRead(d, m)
}


func resourceEC2InstanceDelete(d *schema.ResourceData, m interface{}) error {
    // Load the default AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
    if err != nil {
        return err
    }

    // Create an EC2 client
    svc := ec2.NewFromConfig(cfg)

    // Use the TerminateInstances API to terminate the instance
    _, err = svc.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
        InstanceIds: []string{d.Id()},
    })
    if err != nil {
        return err
    }

    // Remove the instance from the Terraform state
    d.SetId("")

    return nil
}

