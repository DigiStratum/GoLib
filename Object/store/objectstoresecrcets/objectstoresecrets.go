package objectstoresecrets

/*

This is an ObjectStore specifically designed to manage and provide secrets in a secure way.

The store will be initialized with the master secret which gets it access to a mysql ObjectStore,
which gets it the rest of the secrets (which are also encrypted at rest with the secret).
In order to support rotation of the master secret, the stored secrets should be indexed according to
some fingerprint of the secret to indicate which set of credentials are encrypted with it. By
doing this, we can support two copies of the stored secrets: one for the old master secret, and
one for the new master secret, as it is rotated out. This enables applications time to rotate out
the old secret for the new and have both in a working state for some overlap period before the old
is removed completely.

Thinking to use AWS Secrets Manager (https://aws.amazon.com/secrets-manager/) to manage the master
secret. We can make the master secret available at application launch time (likely with environment
variable that can be cleared out or some other authorization scheme that permits the application to
retrieve it at runtime from the Secrets Manager). Once the application has the master secret, it can
instantiate this ObjectStore with configuration pointing to the SecretStore, and use the master
secret to manage the records there.

There must be two modes of access and instantiation: admin and normal.

In admin access mode, we can create, update, and delete secrets in the SecretStore. But in normal
access mode, we only allow reading credentials. Potentially we can use two different keys with some
sort of asymmetrical key cryptography that allows us to create/update secrets with the admin key,
but the normal key can only read those secrets, not create them.

ref: https://medium.com/@Raulgzm/golang-cryptography-rsa-asymmetric-algorithm-e91363a2f7b3

*/
