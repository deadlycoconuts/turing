import React, { useEffect, useMemo, useState } from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiSearchBar,
  EuiSpacer,
  EuiPanel,
  EuiText,
  EuiLink
} from "@elastic/eui";
import { LazyLog, ScrollFollow } from "react-lazylog-with-emitter";
import { slugify } from "@caraml-dev/ui-lib";
import isArray from "lodash/isArray";

export const PodLogsViewer = ({
  components,
  emitter,
  query,
  onQueryChange,
  batchSize,
  podLogUrls,
}) => {
  const filters = useMemo(
    () => [
      {
        type: "field_value_toggle_group",
        field: "component_type",
        items: components,
      },
      {
        type: "field_value_selection",
        field: "tail_lines",
        name: !!query.tail_lines
          ? `Last ${query.tail_lines} records`
          : "From the container start",
        // single select behavior is buggy in recent versions of Eui
        // See: https://github.com/elastic/eui/issues/6271
        multiSelect: "or",
        options: [
          {
            value: "100",
            name: "Last 100 records",
          },
          {
            value: "1000",
            name: "Last 1000 records",
          },
          {
            value: "all",
            name: "From the container start",
          },
        ],
      },
    ],
    [components, query.tail_lines]
  );

  const [filterValues, setFilterValues] = useState({});

  useEffect(() => {
    const filterValues = Object.fromEntries([
      ...Object.entries(query || {}).filter(([k]) =>
        filters.some((f) => f.field === k)
      ),
      ["tail_lines", query.tail_lines || "all"],
    ]);

    setFilterValues((old) =>
      JSON.stringify(old) !== JSON.stringify(filterValues) ? filterValues : old
    );
  }, [query, filters, components, setFilterValues]);

  const searchQuery = useMemo(() => {
    return Object.entries(filterValues)
      .map(([k, v]) => `${k}:"${v}"`)
      .join(" ");
  }, [filterValues]);

  const onChange = ({ query: { ast }, error }) => {
    if (!error) {
      let newFilterValues = {
        ...filterValues,
        ...ast.clauses.reduce((acc, { field, value }) => {
          // This will ensure that, if there are multiple values (as with field_value_selection using OR),
          // only the latest one is selected
          acc[field] = value;
          return acc;
        }, {}),
      };

      if (JSON.stringify(newFilterValues) !== JSON.stringify(filterValues)) {
        if (isArray(newFilterValues.tail_lines)) {
          // Grab the first value, as we are using multi-select with OR
          newFilterValues.tail_lines = newFilterValues.tail_lines[0]
        }
        if (newFilterValues.tail_lines === "all") {
          delete newFilterValues["tail_lines"];
          newFilterValues.head_lines = batchSize;
        }
        onQueryChange(() => newFilterValues);
      }
    }
  };

  const search = {
    query: searchQuery,
    box: {
      readOnly: true,
    },
    filters,
    onChange,
  };

  return (
    <>
      {
        Object.keys(podLogUrls).length !== 0 &&
        (
          <>
            <EuiPanel>
              <EuiFlexGroup direction="row" alignItems="center">
                <EuiFlexItem style={{marginTop:0, marginBottom:0}} grow={false}>
                  <EuiText  style={{ fontSize: '14px', fontWeight:"bold"}}>Pod Logs</EuiText>
                </EuiFlexItem>
                {Object.entries(podLogUrls).map(([component,url])=> (
                  <EuiFlexItem style={{marginTop:0, marginBottom:0, paddingLeft:"10px", textTransform: "capitalize"}} key={component} grow={false}>
                    <EuiText size="xs" >
                      <EuiLink href={url} target="_blank" external>{component.replace(new RegExp("_", "g"), " ")}</EuiLink>
                    </EuiText>
                  </EuiFlexItem>
                ))}
              </EuiFlexGroup>
            </EuiPanel>
            <EuiSpacer size="s" />
          </>
        )
      }
      <EuiPanel>
        <EuiFlexGroup
          direction="column"
          gutterSize="none">
          <EuiFlexItem grow={false}>
            <EuiSearchBar {...search} />
          </EuiFlexItem>
          <EuiFlexItem grow={false}>
            <EuiSpacer size="s" />
          </EuiFlexItem>
          <EuiFlexItem grow={true}>
            <ScrollFollow
              startFollowing={true}
              render={({ onScroll, follow }) => (
                <LazyLog
                  key={slugify(searchQuery)}
                  eventSource={emitter}
                  extraLines={1}
                  onScroll={onScroll}
                  follow={follow}
                  height={640}
                  caseInsensitive
                  enableSearch
                  selectableLines
                />
              )}
            />
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiPanel>
    </>
  );
};
