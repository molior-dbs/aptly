package api

import (
	"strconv"

	"github.com/aptly-dev/aptly/task"
	"github.com/gin-gonic/gin"
)

// GET /tasks
func apiTasksList(c *gin.Context) {
	list := context.TaskList()
	c.JSON(200, list.GetTasks())
}

// POST /tasks/clear
func apiTasksClear(c *gin.Context) {
	list := context.TaskList()
	list.Clear()
	c.JSON(200, gin.H{})
}

// GET /tasks-wait
func apiTasksWait(c *gin.Context) {
	list := context.TaskList()
	list.Wait()
	c.JSON(200, gin.H{})
}

// GET /tasks/:id/wait
func apiTasksWaitForTaskByID(c *gin.Context) {
	list := context.TaskList()
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	task, err := list.WaitForTaskByID(int(id))
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, task)
}

// GET /tasks/:id
func apiTasksShow(c *gin.Context) {
	list := context.TaskList()
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var task task.Task
	task, err = list.GetTaskByID(int(id))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	c.JSON(200, task)
}

// GET /tasks/:id/output
func apiTasksOutputShow(c *gin.Context) {
	list := context.TaskList()
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var output string
	output, err = list.GetTaskOutputByID(int(id))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	c.JSON(200, output)
}

// GET /tasks/:id/detail
func apiTasksDetailShow(c *gin.Context) {
	list := context.TaskList()
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var detail interface{}
	detail, err = list.GetTaskDetailByID(int(id))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	c.JSON(200, detail)
}

// DELETE /tasks/:id
func apiTasksDelete(c *gin.Context) {
	list := context.TaskList()
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var delTask task.Task
	delTask, err = list.DeleteTaskByID(int(id))
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, delTask)
}
