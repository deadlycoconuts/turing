from datetime import date, datetime
from dateutil.tz import tzutc
from typing import Optional, Union, Any, Dict
from turing import generated as client
import turing.ensembler


def json_serializer(o):
    if isinstance(o, (date, datetime)):
        return o.isoformat()
    if isinstance(
        o, (client.model_utils.ModelNormal, client.model_utils.ModelComposed)
    ):
        return o.to_dict()


def utc_date(date_str: str):
    return datetime.strptime(date_str, "%Y-%m-%dT%H:%M:%S.%fZ").replace(tzinfo=tzutc())


class MyTestEnsemblerJob(turing.ensembler.PyFunc):
    import pandas

    def __init__(self, default: float):
        self._default = default

    def initialize(self, artifacts: dict):
        pass

    def ensemble(
        self,
        input: Union[pandas.Series, Dict[str, Any]],
        predictions: Union[pandas.Series, Dict[str, Any]],
        **kwargs,
    ) -> Any:
        if input["treatment"] in predictions:
            return predictions[input["treatment"]]
        else:
            return self._default


class MyTestEnsemblerService(turing.ensembler.PyFunc):
    def __init__(self, default: float):
        self._default = default

    def initialize(self, artifacts: dict):
        pass

    def ensemble(
        self,
        enricher_response,
        **kwargs,
    ) -> Any:
        return enricher_response
