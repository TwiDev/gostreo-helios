import asyncio
import logging
from dataclasses import dataclass
from datetime import timedelta

from temporalio import activity, workflow
from temporalio.client import Client
from temporalio.worker import Worker

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)

@dataclass
class AmazingPayload():
    name: str
    age: int

@activity.defn(name="say-hello-activity")
async def say_hello_activity(payload: AmazingPayload) -> str:
    activity.logger.info("Running activity with parameter %s" % payload)
    return f"Hello, {payload.name}. Today your age is: {payload.age}!"

@workflow.defn(name="say-hello-workflow")
class GreetingWorkflow:
    @workflow.run
    async def run(self, payload: AmazingPayload) -> None:
        result = await workflow.execute_activity(
            say_hello_activity,
            payload,
            start_to_close_timeout=timedelta(seconds=10),
        )
        workflow.logger.info("Result: %s", result)

async def main():
    # Create client to localhost on default namespace
    client = await Client.connect("localhost:7233")
    worker = Worker(
        client,
        task_queue="cron",
        workflows=[GreetingWorkflow],
        activities=[say_hello_activity],
    )
    logger.info("Starting worker...")
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())