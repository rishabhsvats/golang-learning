# Starting EC2 instance using aws-sdk-go-v2

1. Configure `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` 

2. Run main.go whhich will output the private key and Instance ID. The private key can be used to ssh into the EC2 server instance.
   ~~~
   $ go run *.go
   Instance ID: i-0dd83039e0b1cfbdf
   ~~~

3. Use the created private key `go-aws-ec2.pem` to ssh into the EC2 instance
