resource "aws_dynamodb_table" "data_table" {
  name     = "data_table"
  hash_key = "id"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "id"
    type = "S"
  }
}