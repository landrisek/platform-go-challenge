# GlobalWebIndex Engineering Challenge

## Introduction

Hello potential colleagues!

The presented implementation is a stack of simple microservices that demonstrate how the tasks outlined in the acceptance criteria can be achieved while ensuring security, reliability, and ease of extensibility. The approach follows the principle that each microservice should focus on doing one thing well. 

## Local Development

The local development is driven by a Makefile, which means you can run the "make help" command to get descriptions of each task.

The Makefile runs shell scripts that build the Go application and execute the tasks accordingly. The scripts handle the build and execution processes, aiming for a user-friendly experience. Please make sure you have GNU Make installed on your system to use the Makefile.

We have made a deliberate decision to structure our development environment and deployment process in a way that allows for flexibility and efficient development. Instead of incorporating the asset microservice directly into the Docker Compose stack, we have chosen to run it separately from the stack while including all related microservices (user, blacklist) and data layer containers (Redis, MySQL, Vault, migration side-car) within the stack.

This decision is based on the assumption that most of the development work will be focused on the asset microservice, requiring multiple runs and rebuilds during the development process. By separating the asset microservice from the stack, we can avoid the time and resource consumption that would occur if we had to rebuild the entire stack for small changes in a specific microservice.

However, it is important to note that this "cherry-pick" approach can be adapted based on the specific microservice being worked on. For example, if the user service becomes the primary focus, it can be excluded out of the stack while the asset microservice will be added into docker-compose.yml file, allowing developing user service separately.

In a typical scenario, each microservice would have its own repository, allowing for independent development and deployment. Internal packages can also have their own repositories, enabling secure distribution through a private proxy.

By adopting this approach, we can optimize development speed, resource utilization, and collaboration while ensuring that each microservice remains decoupled and independently deployable.

General recommended approach (and base for end-to-end test) is (should be self-explanaible, check make help for more info):
1] make run-dev (wait until finish)
2] make create-user (after fully user microservice standup, can be postpone by provisioning)
3] make run-asset (in separate terminal)
4] make create-asset
5] make read-assets
6] make update-asset
7] make read-asset
8] make delete-asset
9] make read-asset

## Code and comments

In the code, you will find two types of comments:

1] Proper comments that mirror the documentation status and are typically left in the code as part of the documentation. These comments provide explanations and details about the code implementation.

2] HINT. Comments starting with this words I would not typically left in the final codebase, but they serve the purpose of illustrating the evolution of the code and sharing my thoughts during the development process.

3] TODO. These are not necessary for implementation to meet acceptance criteria, however they are descirbing how can we continue

## Data layer

The attached UML graph provides further explanation and illustrates the relationship between the microservices and the vault, which securely stores expirable passwords for MySQL. The decision to use MySQL was based on its simplicity and its well-designed capability to maintain relationships between users and their assets. 
For structuring user table it was itentionaly not used any GDPR data to not overcomplicated implementation, but if they would be, it would be necessary to introduce proper hashing mechanism converting sensitive data into a fixed-length string, making it difficult for unauthorized access. This would provide an additional layer of protection for sensitive information accessed through internal UI tools. 

Redis was chosen as the primary database for purposes such as blacklisting due to its specific features and characteristics that align well with the requirements of this functionality. Redis is an in-memory data store that excels in handling key-value pairs and provides fast read and write operations. Additionally, Redis offers features like expiration of keys, which can be leveraged to automatically remove outdated or unused entries from the blacklist. This helps maintain the integrity and relevance of the forbidden word database.

However, it's important to note that while Redis is a performant and efficient database, proper backup and monitoring mechanisms should still be introduced to ensure data reliability and availability. The UML diagram provides insights into the suggested backup and monitoring setup, highlighting the need for regular backups and monitoring of Redis instances to prevent data loss and ensure system stability.

## Cache and hot spots

While the implementation of Redis cache was not included in the current version of the system, it is a consideration for future enhancements. However, it is important to note that we prioritize other aspects such as clean design, proper testing, and overall system reliability over the implementation of a caching layer.

By focusing on clean design principles and ensuring comprehensive test coverage, we aim to build a robust and scalable microservices architecture. These aspects lay a solid foundation for the system's performance and maintainability.

While Redis caching can bring performance benefits by reducing database load and improving response times, it is considered a secondary optimization compared to the core functionality and quality of the system. As the system evolves, we will evaluate the need for implementing Redis cache based on performance benchmarks and user requirements.

## Configuration
For configuration was used environmental variables, for localhost development save under artifacts/.env and artifacts/.asset.env as a a lightweight and flexible option for simple configurations. For more complex structure can be introduce loading through yaml files, thus both needs to be properly secured.

## User Microservice

In this particular case, the requirements specified CRUD operations on user assets, rather than CRUD operations on the user entity itself. However, it is evident that we cannot create an asset without an associated user.
To fulfill this purpose, a simple microservice called "user" was created. The principle of each microservice focusing on doing one thing well is crucial, especially in cases where one microservice is down. By following this principle, the functionality of other microservices remains unaffected even if a particular microservice, such as the one handling assets, is temporarily unavailable. This allows for better fault isolation and ensures that the overall system can continue to operate without disruption. It also simplifies maintenance and updates, as changes to one microservice are less likely to impact others. This modular and decoupled approach enhances the resilience and reliability of the system as a whole.

## Asset Microservice

In the code dedicated to the microservice handling CRUD operations on assets, the Saga pattern was chosen. This pattern is well-suited for microservice architecture as it allows for the execution of various nested operations while keeping the codebase as organized as possible. Although the Saga pattern tends to be sequential by default, we have incorporated self-thread concurrency patterns at a lower level to improve performance. These patterns are managed by pooling goroutines, ensuring that the concurrency remains within reasonable boundaries and enhancing overall efficiency.

## Blacklist Microservice

The blacklist microservice was implemented as a simple example to showcase the effectiveness of the saga pattern in handling multiple operations. Its specific purpose is to replace forbidden words with appropriate alternatives. While it may seem like a straightforward task, integrating it into the saga pattern demonstrates the flexibility and scalability of this architectural approach.

By incorporating the blacklist microservice within the saga, we can ensure that any content or input containing prohibited words is automatically modified before further processing. This not only helps maintain the integrity and appropriateness of the data but also allows for easy customization and expansion of the blacklist functionality.

Moreover, if we will come to decision, that blacklist should be replaced by completly new microservice (e.g. in different language) or it will be replaced by better, external (even paid) service, it will be very easy to do so.

## Backup Microservice

Was not implemented. Propose design was to archive soft-delete entries to something like glacier or separate instance of redis and hard-delete them. 

## Security

For authentication was used simple token stored both in redis and mysql. Crud permissions are handled by table permissions differ each operation separately for each token. Additionaly, can be introduce table permissions_users which will be allowing only limited number of user withing given security group defined by token.

## Tests

Last, but not least.