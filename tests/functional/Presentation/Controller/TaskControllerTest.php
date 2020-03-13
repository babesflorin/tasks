<?php

namespace App\Tests\Functional\Presentation\Controller;

use App\Presentation\Controller\TaskController;
use App\Tests\DataFixtures\TasksFixtures;
use Liip\TestFixturesBundle\Test\FixturesTrait;
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpFoundation\Response;

class TaskControllerTest extends WebTestCase
{
    use FixturesTrait;

    public function setUp(): void
    {
        parent::setUp();
        $this->loadFixtures([TasksFixtures::class]);
    }

    public function testGetTask()
    {
        $client = $this->createClient();
        $taskId = 1;
        $client->request("GET", "/api/task/$taskId");
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_OK, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertSame('Task name 0', $response['data']['name']);
        $this->assertSame('Task description 0', $response['data']['description']);
        $this->assertSame((new \DateTime())->format('Y-m-d'), $response['data']['when']);
        $this->assertSame(false, $response['data']['done']);
    }

    public function testGetTaskNotFound()
    {
        $client = $this->createClient();
        $taskId = 9999999;
        $client->request("GET", "/api/task/$taskId");
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_NOT_FOUND, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertEmpty($response['data']);
        $this->assertSame('Task not found!', $response['error']);
    }


    public function testAddTask()
    {
        $client = $this->createClient();
        $taskId = 1;
        $taskData = [
            'name' => "Task name 6",
            'description' => "Task description 6",
            'when' => (new \DateTime())->modify('+5 day')->format('Y-m-d'),
        ];
        $client->request(
            "POST",
            "/api/task",
            [],
            [],
            [],
            json_encode($taskData)
        );
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_OK, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertSame($taskData['name'], $response['data']['name']);
        $this->assertSame($taskData['description'], $response['data']['description']);
        $this->assertSame($taskData['when'], $response['data']['when']);
        $this->assertSame(false, $response['data']['done']);
    }

    public function testUpdateTask()
    {
        $client = $this->createClient();
        $taskData = [
            'id' => 2,
            'name' => "Task name 22",
            'description' => "Task description 44",
            'when' => (new \DateTime())->modify('+2 day')->format('Y-m-d'),
        ];
        $client->request(
            "PUT",
            "/api/task",
            [],
            [],
            [],
            json_encode($taskData)
        );
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_OK, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertSame($taskData['id'], $response['data']['id']);
        $this->assertSame($taskData['name'], $response['data']['name']);
        $this->assertSame($taskData['description'], $response['data']['description']);
        $this->assertSame($taskData['when'], $response['data']['when']);
        $this->assertSame(false, $response['data']['done']);
    }

    public function testUpdateTaskNotFound()
    {
        $client = $this->createClient();
        $taskData = [
            'id' => 9999999,
            'name' => "Task name 22",
            'description' => "Task description 44",
            'when' => (new \DateTime())->modify('+2 day')->format('Y-m-d'),
        ];
        $client->request(
            "PUT",
            "/api/task",
            [],
            [],
            [],
            json_encode($taskData)
        );
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_NOT_FOUND, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertEmpty($response['data']);
        $this->assertSame('Task not found!', $response['error']);
    }

    public function testUpdateTaskInvalidData()
    {
        $client = $this->createClient();
        $taskData = [
        ];
        $client->request(
            "PUT",
            "/api/task",
            [],
            [],
            [],
            json_encode($taskData)
        );
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_BAD_REQUEST, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertEmpty($response['data']);
        $this->assertSame('Task is not valid!', $response['error']);
        $messages = [
            'Task name is not valid!',
            'We need an id to know which entity to update!',
            'Task description is not valid!',
            'Task must have a date!',
        ];
        $this->assertSame($messages, $response['messages']);
    }

    public function testUpdateInvalidJson()
    {
        $client = $this->createClient();
        $client->request(
            "PUT",
            "/api/task",
            [],
            [],
            [],
            "/"
        );
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_BAD_REQUEST, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertEmpty($response['data']);
        $this->assertSame('Request must be json!', $response['error']);
    }

    public function testGetTasks()
    {
        $client = $this->createClient();
        $taskId = 1;
        $client->request("GET", "/api/task");
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_OK, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertNotEmpty($response['data']);
        $this->assertCount(5, $response['data']);
    }

    public function testCompleteTask()
    {
        $client = $this->createClient();
        $taskId = 1;
        $client->request("PUT", "/api/task/$taskId/complete");
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_OK, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertSame('Task name 0', $response['data']['name']);
        $this->assertSame('Task description 0', $response['data']['description']);
        $this->assertSame((new \DateTime())->format('Y-m-d'), $response['data']['when']);
        $this->assertSame(true, $response['data']['done']);
    }

    public function testCompleteTaskNotFound()
    {
        $client = $this->createClient();
        $taskId = 9999999;
        $client->request("PUT", "/api/task/$taskId/complete");
        $this->assertTrue(
            $client->getResponse()->headers->contains(
                'Content-Type',
                'application/json'
            )
        );
        $this->assertSame(Response::HTTP_NOT_FOUND, $client->getResponse()->getStatusCode());
        $response = json_decode($client->getResponse()->getContent(), true);
        $this->assertSame(JSON_ERROR_NONE, json_last_error());
        $this->assertEmpty($response['data']);
        $this->assertSame('Task not found!', $response['error']);
    }
}
