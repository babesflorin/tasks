<?php

namespace App\Tests\Domain\Service;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Entity\TaskCollection;
use App\Domain\Repository\TaskRepositoryInterface;
use App\Domain\Service\TaskService;
use App\Domain\Transformer\TaskCollectionDtoToCollectionTransformer;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Domain\Validator\TaskValidator;
use App\Tests\TestCase;

class TaskServiceTest extends TestCase
{
    private $repositoryMock;

    private $validatorMock;

    private $transformerMock;

    private $collectionTransformerMock;

    protected function setUp()
    {
        parent::setUp();
        $this->repositoryMock = $this->getMock(
            TaskRepositoryInterface::class,
            ['saveTask', 'getTasks', 'findTaskById']
        );
        $this->validatorMock = $this->getMock(TaskValidator::class, ['validate']);
        $this->transformerMock = $this->getMock(TaskDtoToEntityTransformer::class, ['transform', 'reverseTransform']);
        $this->collectionTransformerMock = $this->getMock(
            TaskCollectionDtoToCollectionTransformer::class,
            ['transform', 'reverseTransform']
        );
    }

    public function testAddTask()
    {
        $taskMock = $this->getMock(Task::class);
        $taskDto = new TaskDto();
        $this->repositoryMock->expects(self::once())->method('saveTask')->with($taskMock)->willReturn($taskMock);
        $this->validatorMock->expects(self::once())->method('validate')->with($taskDto)->willReturn(true);
        $this->transformerMock->expects(self::once())->method('transform')->with($taskDto)->willReturn($taskMock);
        $this->transformerMock->expects(self::once())->method('reverseTransform')->with($taskMock)->willReturn(
            $taskDto
        );
        $service = new TaskService(
            $this->repositoryMock,
            $this->validatorMock,
            $this->transformerMock,
            $this->collectionTransformerMock
        );
        $this->assertSame($taskDto, $service->addTask($taskDto));
    }

    public function testGetAllTasks()
    {
        $taskCollection = new TaskCollection();
        $taskCollectionDto = new TaskCollectionDto();
        $this->repositoryMock->expects(self::once())->method('getTasks')->willReturn($taskCollection);
        $this->collectionTransformerMock->expects(self::once())->method('reverseTransform')->with($taskCollection)->willReturn(
            $taskCollectionDto
        );
        $service = new TaskService(
            $this->repositoryMock,
            $this->validatorMock,
            $this->transformerMock,
            $this->collectionTransformerMock
        );

        $this->assertSame($taskCollectionDto, $service->getAllTasks());
    }
}
