<?php

namespace App\Tests\Unit\Domain\Service;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Entity\TaskCollection;
use App\Domain\Exception\CouldNotDeleteException;
use App\Domain\Exception\TaskNotFoundException;
use App\Domain\Repository\TaskRepositoryInterface;
use App\Domain\Service\TaskService;
use App\Domain\Transformer\TaskCollectionDtoToCollectionTransformer;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Domain\Validator\TaskValidator;
use App\Tests\Unit\TestCase;

class TaskServiceTest extends TestCase
{
    private $repositoryMock;

    private $validatorMock;

    private $transformerMock;

    private $collectionTransformerMock;
    private $taskService;

    protected function setUp(): void
    {
        parent::setUp();
        $this->repositoryMock = $this->getMock(
            TaskRepositoryInterface::class,
            ['saveTask', 'getTasks', 'findTaskById', 'deleteTask']
        );
        $this->validatorMock = $this->getMock(TaskValidator::class, ['validate']);
        $this->transformerMock = $this->getMock(TaskDtoToEntityTransformer::class, ['transform', 'reverseTransform']);
        $this->collectionTransformerMock = $this->getMock(
            TaskCollectionDtoToCollectionTransformer::class,
            ['transform', 'reverseTransform']
        );
        $this->taskService = new TaskService(
            $this->repositoryMock,
            $this->validatorMock,
            $this->transformerMock,
            $this->collectionTransformerMock
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
        ;
        $this->assertSame($taskDto, $this->taskService->addTask($taskDto));
    }

    public function testGetAllTasks()
    {
        $taskCollection = new TaskCollection();
        $taskCollectionDto = new TaskCollectionDto();
        $this->repositoryMock->expects(self::once())->method('getTasks')->willReturn($taskCollection);
        $this->collectionTransformerMock->expects(self::once())
                                        ->method('reverseTransform')
                                        ->with($taskCollection)
                                        ->willReturn(
                                            $taskCollectionDto
                                        );
        $this->assertSame($taskCollectionDto, $this->taskService->getAllTasks());
    }

    public function testCompleteTask()
    {
        $taskId = 1;
        $taskMock = $this->getMock(Task::class, ['complete']);
        $taskMock->expects(self::once())->method('complete');
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskId)->willReturn($taskMock);
        $this->repositoryMock->expects(self::once())->method('saveTask')->with($taskMock);
        $taskDto = new TaskDto();
        $this->transformerMock->expects(self::once())->method('reverseTransform')->with($taskMock)->willReturn(
            $taskDto
        );
        $this->assertSame($taskDto, $this->taskService->completeTask($taskId));
    }

    public function testCompleteTaskNotFound()
    {
        $taskId = 1;
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskId)->willReturn(null);
        $this->expectException(TaskNotFoundException::class);
        $this->taskService->completeTask($taskId);
    }

    public function testGetTask()
    {
        $taskId = 1;
        $taskMock = $this->getMock(Task::class);
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskId)->willReturn($taskMock);
        $taskDto = new TaskDto();
        $this->transformerMock->expects(self::once())->method('reverseTransform')->with($taskMock)->willReturn(
            $taskDto
        );
        $this->assertSame($taskDto, $this->taskService->getTask($taskId));
    }

    public function testGetTaskNotFound()
    {
        $taskId = 1;
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskId)->willReturn(null);
        $this->expectException(TaskNotFoundException::class);
        $this->taskService->getTask($taskId);
    }

    public function testUpdateTask()
    {
        $taskDto = new TaskDto();
        $taskDto->id = 1;
        $this->validatorMock->expects(self::once())->method('validate')->with($taskDto);
        $taskMock = $this->getMock(Task::class);
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskDto->id)->willReturn($taskMock);
        $this->repositoryMock->expects(self::once())->method('saveTask')->with($taskMock);
        $this->transformerMock->expects(self::once())->method('transform')->with($taskDto, $taskMock)->willReturn(
            $taskMock
        );
        $this->transformerMock->expects(self::once())->method('reverseTransform')->with($taskMock)->willReturn(
            $taskDto
        );
        $this->assertSame($taskDto, $this->taskService->updateTask($taskDto));
    }

    public function testDeleteTask()
    {
        $taskId = 5;
        $taskMock = $this->getMock(Task::class);
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskId)->willReturn($taskMock);
        $this->repositoryMock->expects(self::once())->method('deleteTask')->with($taskMock)->willReturn(true);
        $taskDto = new TaskDto();
        $this->transformerMock->expects(self::once())->method('reverseTransform')->with($taskMock)->willReturn(
            $taskDto
        );
        $this->assertSame($taskDto, $this->taskService->deleteTask($taskId));
    }

    public function testDeleteTaskNotFound()
    {
        $taskId = 5;
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskId)->willReturn(null);
        $this->expectException(TaskNotFoundException::class);
        $this->taskService->deleteTask($taskId);
    }


    public function testDeleteTaskRepositoryException()
    {
        $taskId = 5;
        $taskMock = $this->getMock(Task::class);
        $this->repositoryMock->expects(self::once())->method('findTaskById')->with($taskId)->willReturn($taskMock);
        $this->repositoryMock->expects(self::once())->method('deleteTask')->with($taskMock)->willThrowException(
            new \Exception()
        );
        $this->expectException(CouldNotDeleteException::class);
        $this->taskService->deleteTask($taskId);
    }
}
