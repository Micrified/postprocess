import pandas as pd
import plotly.graph_objs as go
from plotly.offline import init_notebook_mode, iplot
# init_notebook_mode()

df = pd.read_csv('filtered_data.csv')
df.columns = ['executor', 'chain', 'callback', 'start', 'duration']
# df['duration'] = df['end'] - df['start']

fig = go.Figure(
    layout = {
        'title': {'text': 'Task timeplot'},
        'xaxis_title': {'text': 'Time (us)'},
        'yaxis_title': {'text': 'Chains'},
        'barmode': 'stack',
        'xaxis': {'automargin': True},
        'yaxis': {'automargin': True}}#, 'categoryorder': 'category ascending'}}
)

for callback, callback_df in df.groupby('callback'):
    fig.add_bar(x=callback_df.duration,
                y=callback_df.chain,
                base=callback_df.start,
                orientation='h',
                showlegend=False,
                name=callback)


fig.update_layout(
    yaxis = dict(
        tickmode = 'linear',
        tick0 = 0,
        dtick = 1
    )
    # ,
    # xaxis = dict(
    #     tickmode = 'linear',
    #     tick0 = 0,
    #     dtick = 500000
    # )
)


# fig.add_annotation(x=358879, y=1,
#             text="170450ns",
#             showarrow=True,
#             arrowhead=1)

# fig.add_shape(type="line",
#     xref="x", yref="y",
#     x0=111686, y0=0, x1=111686, y1=1,
#     line=dict(
#         color="Black",
#         width=3,
#         dash="dot",
#     ),
# )
# fig.update_shapes(dict(xref='x', yref='y'))


# iplot(fig)
fig.show()