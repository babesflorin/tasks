<?php

namespace App\Tests\Unit\Domain\Transformer;

use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Tests\Unit\TestCase;

class TaskDtoToEntityTransfomerTest extends TestCase
{
    /**
     * @var TaskDtoToEntityTransformer
     */
    private $transformer;

    protected function setUp(): void
    {
        parent::setUp();
        $this->transformer = new TaskDtoToEntityTransformer();
    }

    public function testTransform()
    {
        $taskDto = new TaskDto();
        $taskDto->name = "test";
        $taskDto->description = "test";
        $time = new \DateTime();
        $time->setTime(0, 0, 0);
        $taskDto->when = $time->format('Y-m-d');

        $expectedEntity = new Task($taskDto->name, $taskDto->description, $time);
        $transformedDto = $this->transformer->transform($taskDto);
        $this->assertEquals($expectedEntity->getName(), $transformedDto->getName());
        $this->assertEquals($expectedEntity->getDescription(), $transformedDto->getDescription());
        $this->assertEquals($expectedEntity->getWhen(), $transformedDto->getWhen());
    }

    public function testTransformExistingEntity()
    {
        $taskDto = new TaskDto();
        $taskDto->name = "test";
        $taskDto->description = "test";
        $time = new \DateTime();
        $time->setTime(0, 0, 0);
        $taskDto->when = $time->format('Y-m-d');

        $entityToUpdate = new Task("not the dto name", "not the dto description", $time);
        $expectedEntity = new Task($taskDto->name, $taskDto->description, $time);
        $transformedDto = $this->transformer->transform($taskDto, $entityToUpdate);
        $this->assertEquals($expectedEntity->getName(), $transformedDto->getName());
        $this->assertEquals($expectedEntity->getDescription(), $transformedDto->getDescription());
        $this->assertEquals($expectedEntity->getWhen(), $transformedDto->getWhen());
    }

    public function testReverseTransform()
    {
        $expectedDto = new TaskDto();
        $expectedDto->name = "test";
        $expectedDto->description = "test";
        $time = new \DateTime();
        $expectedDto->when = $time->format('Y-m-d');
        $expectedDto->done = false;
        $entity = new Task($expectedDto->name, $expectedDto->description, $time);
        $expectedDto->createdAt = $entity->getCreatedAt();
        $expectedDto->updatedAt = $entity->getUpdateAt();
        $this->assertEquals($expectedDto, $this->transformer->reverseTransform($entity));
    }
}
