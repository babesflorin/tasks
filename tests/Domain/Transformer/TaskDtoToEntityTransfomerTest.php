<?php

namespace App\Tests\Domain\Transformer;

use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Tests\TestCase;

class TaskDtoToEntityTransfomerTest extends TestCase
{
    /**
     * @var TaskDtoToEntityTransformer
     */
    private $transformer;

    protected function setUp()
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
        $time->setTime(0,0,0);
        $taskDto->when = $time->format('Y-m-d');

        $expectedEntity = new Task($taskDto->name, $taskDto->description, $time);
        $this->assertEquals($expectedEntity, $this->transformer->transform($taskDto));
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
        $this->assertEquals($expectedDto, $this->transformer->reverseTransform($entity));
    }
}
