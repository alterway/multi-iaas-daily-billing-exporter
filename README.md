# Multi-iaas-billing exporter
 
Multi-iaas-billing-exporter enables to collect, unify and expose daily billing from AWS and GCP providers. The aim is to trigger the daily billing from AWS and GCP and expose it to prometheus to provide a global vision about your resources billing in multiple providers. The main contributions are:

- One exporter for multiple iaas providers
- Unification of collected billing information (day, service, cost, currency)
- Centralization of billing in prometheus

## Requirements:

- Kubernetes 
- A prometheus instance
- AWS and GCP accounts with a full access right to billing 

**Note that** it is possible to modify the `app/controller.go` to run only one of GCP or AWS billing.


## Multi-iaas-billing-exporter architecture

In the following architecture, we shwocase how multi-iaas-billing-exporter exports billing from GCP and AWS and send it to prometheus. The resulted metrics `aws_cost` and `gcp_cost` contain the following information:

- SERVICE: the AWS or GCP service name
- DAY : the day of the billing
- CURRENCY: the currency (e.g. euros)
- COST: the cost of this service during that DAY

![prom](/img/multi-iaas.png)

## USAGE:

To deploy the multi-iaas-billing-exporter, you should enable billing and configure your secret key access to both GCP and AWS.

### Enabling billing data export:
#### 1. Google Cloud Platform (GCP): 

In order to export billing from GCP, we consider BiQuery which is a Google Cloud's fully managed data warehouse that lets you export and query your GCP billing. Cloud Billing export to BigQuery enables you to export detailed Google Cloud billing data (such as usage, cost estimates, and pricing data) automatically throughout the day to a BigQuery dataset that you specify. In our case, we query billing information from the dataset you create in BigQuery. 

To export billing data into a BigQuery dataset follow the steps in this documentation: [BigQuerry](https://cloud.google.com/billing/docs/how-to/export-data-bigquery) 

Once your BigQuery dataset created and billing export enabled, add your `GCP_PROJECTID` and `GCP_TABLE` in the `.env` file.

#### 2. Amazon Web Services (AWS): 

To export AWS billing, we consider the AWS API. To activate your acces to AWS billing and cost management follow the steps described 
in [AWS docs](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/control-access-billing.html)

- 
### AWS and GCP key access configuration 

- To configure the GCP access key:
    - First create and download your key [here](https://cloud.google.com/docs/authentication/getting-started) in a json file that you name `key.json`.
    - Then:

    ```
    kubectl create secret generic gcp-key --from-file=/path/key.json
    ```

- To configure the AWS access key:
    - First create and download your AWS account credentials and configuration files [see](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).
    - Then:
    ```
    k create secret generic aws-key --from-file=/path/credentials --from-file=/path/config
    ```


### To run: 

`kubectl apply -f deploy/`

## Results:

In your prometheus instance, you can find `aws_cost` and `gcp_cost` metrics after the intervall specified in `manifest.yaml`. This interval represents the frequency of daily billing query.  ** Note that** a very small intervall may imply additional cost while querying billing from AWS and GCP providers. 
The following screenshot showcases the ``gcp_cost` metric in prometheus.
![prom](/img/result.png)


Made with **<3** by the **R&D team @ [`Alter way`](https://www.alterway.fr/)**

