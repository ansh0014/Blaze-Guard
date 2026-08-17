from langchain.agents import initialize_agent, AgentType
from agents.llm import llm
from agents.tools import(
    fire_analysis_tool,
    incident_report_tool,
    citizen_alert_tool
)

tools = [
    fire_analysis_tool,
    incident_report_tool,
    citizen_alert_tool
]

agent_executor = initialize_agent(
    tools,
    llm,
    agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True
)