package poll

func (pr *PollRepository) DeletePoll(pollID uint64, userID uint64) error {
	const context = "PollRepository.DeletePoll"

	resp, err := pr.tarantoolUtils.ProcessCall(pr.Conn, "delete_poll", []interface{}{pollID, userID}, context)
	if err != nil {
		return err
	}

	result, err := pr.tarantoolUtils.ExtractMap(resp, context)
	if err != nil {
		return err
	}

	if err := pr.tarantoolUtils.HandleTarantoolError(result, context); err != nil {
		return err
	}

	return nil
}
