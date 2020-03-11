<?php

namespace App\Presentation\Controller;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Entity\TaskCollection;
use App\Domain\Service\TaskService;
use App\Presentation\Transformer\TaskTransformer;
use League\Fractal\Manager;
use League\Fractal\Resource\Collection;
use League\Fractal\Resource\Item;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\ParamConverter;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\Routing\Annotation\Route;
use Swagger\Annotations as SWG;

/**
 * @Route("/api/task")
 * @SWG\Tag(name="tasks")
 */
class TaskController extends AbstractController
{
    /**
     * @var TaskService
     */
    private $service;

    /**
     * @var Manager
     */
    private $fractal;

    /**
     * @var TaskTransformer
     */
    private $taskTransformer;

    public function __construct(TaskService $service, Manager $fractal, TaskTransformer $taskTransformer)
    {
        $this->service = $service;
        $this->fractal = $fractal;
        $this->taskTransformer = $taskTransformer;
    }

    /**
     * @Route(
     *     "",
     *     name="task_add",
     *     methods={"POST"},
     *     format="application/json",
     *     requirements={
     *          "_format" : "application/json"
     *      }
     * )
     * @SWG\Post(
     *      @SWG\Parameter(
     *          name="body",
     *          in="body",
     *          format="application/json",
     *          @SWG\Schema(ref="#/definitions/TaskRequest")
     *      ),
     *      @SWG\Response(
     *          response=200,
     *          description="Returns added task",
     *          @SWG\Schema(ref="#/definitions/TaskResponse")
     *      )
     * )
     * @ParamConverter("task", class=TaskDto::class)
     */
    public function addTask(TaskDto $task)
    {
        $taskDto = $this->service->addTask($task);
        $resource = new Item($taskDto, $this->taskTransformer);

        return $this->json(
            $this->fractal->createData($resource)->toArray()
        );
    }

    /**
     * @Route(
     *     "",
     *     name="tasks_fet",
     *     methods={"GET"},
     *     format="application/json",
     *     requirements={
     *          "_format" : "application/json"
     *      }
     * )
     * @SWG\Response(
     *     response=200,
     *     description="Returns the tasks",
     *     @SWG\Schema(
     *         type="array",
     *         @SWG\Items(ref="#/definitions/TaskResponse")
     *     )
     * )
     */
    public function getTasks()
    {
        $tasksCollectionDto = $this->service->getAllTasks();
        $resource = new Collection($tasksCollectionDto, $this->taskTransformer);

        return $this->json(
            $this->fractal->createData($resource)->toArray()
        );
    }

    /**
     * @Route(
     *     "/{taskId}/complete",
     *     name="task_complete",
     *     methods={"PUT"},
     *     format="application/json",
     *     requirements={
     *          "_format" : "application/json"
     *      }
     * )
     * @SWG\Put(
     *      @SWG\Parameter(parameter="taskId", name="taskId", type="integer", in="path"),
     *      @SWG\Response(
     *          response=200,
     *          description="Returns the completed task",
     *          @SWG\Schema(ref="#/definitions/TaskResponse")
     *      )
     * )
     */
    public function completeTask(int $taskId)
    {
        $taskDto = $this->service->completeTask($taskId);
        $resource = new Item($taskDto, $this->taskTransformer);

        return $this->json(
            $this->fractal->createData($resource)->toArray()
        );
    }
}
