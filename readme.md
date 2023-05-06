

Example

```json
{
  "hostname": "dummy.rtp.raleigh.ibm.com",
  "whenPattern": "MM/dd/yyyy HH:mm:ss",
  "metadata": {
    "mdata1": "value 1",
    "mdata2": "value 2",
    "mdata3": "completely arbitrary data goes here"
  },
  "types": [
    {
      "id": "type1",
      "name": "type 1",
      "fields": [ "field11", "field12", "field13" ]
    },
    { 
      "id": "type2",
      "name": "type 2",
      "fields": [ "field21", "field22", "field23" ]
    },
    {
      "id": "type3",
      "name": "type 3",
      "fields" :[ "field31", "field32", "field33" ]
    },
    {
      "id": "type4",
      "name": "type 4",
      "subtypes": [ "sub41", "sub42", "sub43"],
      "fields" :[ "field41", "field42", "field43" ]
    },
    {
      "id": "type5",
      "name": "type 5",
      "subtypes": [ "sub51", "sub52", "sub53"],
      "fields" :[ "field51", "field52", "field53" ]
    }
  ],
  "data": [
    { "when": "10/15/2012 10:00:00",
      "type1": [ 1, 2, 3 ],
      "type2": [ 4, 5, 6 ],
      "type3": [ 1, 2 , 3 ],
      "type4": {
                 "sub41": [ 4, 5 , 6 ],
                 "sub42": [ 6, 7 , 8 ],
                 "sub43": [ 8, 9 , 10 ]
      },
      "type5": {
                 "sub51": [ 3, 2 , 1 ],
                 "sub52": [ 4, 3 , 2 ],
                 "sub53": [ 5, 4 , 3 ]
      }
    },
    { "when": "10/15/2012 10:00:30",
      "type1": [ 2, 3 , 4 ],
      "type2": [ 5, 6 , 7 ],
      "type3": [ 1, 2 , 3 ],
      "type4": {
                 "sub41": [ 4, 5 , 6 ],
                 "sub42": [ 6, 7 , 8 ],
                 "sub43": [ 8, 9 , 10 ]
      },
      "type5": {
                 "sub51": [ 3, 2 , 1 ],
                 "sub52": [ 4, 3 , 2 ],
                 "sub53": [ 5, 4 , 3 ]
      }
    },
    { "when": "10/15/2012 10:01:00",
      "type1": [ 1, 2 , 3 ],
      "type2": [ 4, 5 , 6 ],
      "type3": [ 1, 2 , 3 ],
      "type4": {
                 "sub41": [ 4, 5 , 6 ],
                 "sub42": [ 6, 7 , 8 ],
                 "sub43": [ 8, 9 , 10 ]
      },
      "type5": {
                 "sub51": [ 3, 2 , 1 ],
                 "sub52": [ 4, 3 , 2 ],
                 "sub53": [ 5, 4 , 3 ]
      }
    },
    { "when": "10/15/2012 10:01:30",
      "type1": [ 2, 3 , 4 ],
      "type2": [ 5, 6 , 7 ],
      "type3": [ 1, 2 , 3 ],
      "type4": {
                 "sub41": [ 4, 5 , 6 ],
                 "sub42": [ 6, 7 , 8 ],
                 "sub43": [ 8, 9 , 10 ]
      },
      "type5": {
                 "sub51": [ 3, 2 , 1 ],
                 "sub52": [ 4, 3 , 2 ],
                 "sub53": [ 5, 4 , 3 ]
      }
    },
    { "when" : "10/15/2012 10:02:00",
      "type1": [ 1, 2 , 3 ],
      "type2": [ 4, 5 , 6 ],
      "type3": [ 1, 2 , 3 ],
      "type4": {
                 "sub41": [ 4, 5 , 6 ],
                 "sub42": [ 6, 7 , 8 ],
                 "sub43": [ 8, 9 , 10 ]
      },
      "type5": {
                 "sub51": [ 3, 2 , 1 ],
                 "sub52": [ 4, 3 , 2 ],
                 "sub53": [ 5, 4 , 3 ]
      }
    }
  ]
}
```
