<?php

namespace App\Domain\Entity;

class Task
{
    private $id;

    private $name;

    private $description;

    private $when;

    private $done = false;

    private $createdAt;

    private $updatedAt;

    public function __construct(string $name, string $description, \DateTimeInterface $when)
    {
        $this->name = $name;
        $this->description = $description;
        $this->when = $when;
    }

    public function getCreatedAt(): ?\DateTime
    {
        return $this->createdAt;
    }

    public function getUpdateAt(): ?\DateTime
    {
        return $this->updatedAt;
    }

    public function getName(): ?string
    {
        return $this->name;
    }

    public function getDescription(): ?string
    {
        return $this->description;
    }

    public function getWhen(): ?\DateTimeInterface
    {
        return $this->when;
    }

    public function isDone(): ?bool
    {
        return $this->done;
    }

    public function completeTask(): self
    {
        $this->done = true;

        return $this;
    }

    public function getId(): ?int
    {
        return $this->id;
    }

    public function updateTimestamps(): void
    {
        if (null === $this->createdAt) {
            $this->createdAt = new \DateTime();
        }
        $this->updatedAt = new \DateTime();
    }
}
