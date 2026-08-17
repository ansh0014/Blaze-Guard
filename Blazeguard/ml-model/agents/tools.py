import json
from langchain.tools import Tool
from agents.llm import llm
from pipeline import analyze_fire

def fire_analysis(image_path:str):
    result = analyze_fire(image_path)
    return json.dumps(result,indent=2)

fire_analysis_tool = Tool(
    name = "FireAnalysisTool",
    func = fire_analysis,
    description=(
        "Use this tool to analyze a fire image. "
        "Input should be an image file path. "
        "Returns structured fire intelligence including severity, radius, routing."
    )
)

def generate_incident_report(fire_data: str):

    prompt = f"""
    You are a wildfire command center AI.

    Generate a structured incident report from the following fire analysis data:

    {fire_data}

    Include:
    - Fire Classification
    - Risk Level
    - Estimated Spread
    - Nearest Fire Station
    - Operational Recommendations
    """

    response = llm.invoke(prompt)
    return response.content

incident_report_tool = Tool(
    name = "IncidentReportTool",
    func = generate_incident_report,
    description = (
        "Use this tool to generate a professional wildfire incident report. "
        "Input should be structured fire analysis data in JSON format."
    )
)

def generate_citizen_alert(fire_data: str):

    prompt = f"""
    Create a public safety alert message for citizens based on:

    {fire_data}

    The message must:
    - Be calm and clear
    - Include evacuation guidance if necessary
    - Mention emergency contact advice
    - Avoid technical language
    """

    response = llm.invoke(prompt)
    return response.content

citizen_alert_tool = Tool(
    name = "CitizenAlertTool",
    func = generate_citizen_alert,
    description = (
        "Use this tool to create a public safety alert message "
        "based on fire analysis JSON data."
    )
)