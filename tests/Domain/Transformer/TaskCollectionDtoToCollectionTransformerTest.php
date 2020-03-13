<?php

namespace App\Tests\Domain\Transformer;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Entity\TaskCollection;
use App\Domain\Transformer\TaskCollectionDtoToCollectionTransformer;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Tests\TestCase;

class TaskCollectionDtoToCollectionTransformerTest extends TestCase
{
    /**
     * @var TaskDtoToEntityTransformer
     */
    private $transformer;
    /**
     * @var TaskCollectionDtoToCollectionTransformer
     */
    private $dtoTransformer;

    protected function setUp(): void
    {
        parent::setUp();
        $this->transformer = $this->getMock(TaskDtoToEntityTransformer::class, ['transform', 'reverseTransform']);
        $this->dtoTransformer = new TaskCollectionDtoToCollectionTransformer($this->transformer);
    }

    public function testTransform()
    {
        $taskDto = new TaskDto();
        $taskCollectionDto = new TaskCollectionDto([$taskDto]);
        $task = $this->getMock(Task::class);
        $this->transformer->expects(self::once())->method('transform')->with($taskDto)->willReturn($task);
        $expectCollection = new TaskCollection([$task]);
        $this->assertEquals($expectCollection, $this->dtoTransformer->transform($taskCollectionDto));
    }

    public function testReverseTransform()
    {
        $taskDto = new TaskDto();
        $expectCollection = new TaskCollectionDto([$taskDto]);
        $task = $this->getMock(Task::class);
        $this->transformer->expects(self::once())->method('reverseTransform')->with($task)->willReturn($taskDto);
        $taskCollection = new TaskCollection([$task]);
        $this->assertEquals($expectCollection, $this->dtoTransformer->reverseTransform($taskCollection));
    }
}
