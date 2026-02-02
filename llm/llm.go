package llm

type LLM_API interface {
    Call(prompt string) (response string, err error)
}
