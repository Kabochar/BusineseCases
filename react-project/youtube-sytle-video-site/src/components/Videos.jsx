import { Box, Stack } from "@mui/material";
import React from 'react';
import { ChannelCard, VideoCard } from './';


function Videos({ videos,  direction}) {
    if(!videos?.length) return "loding...";
    return (
        <Stack
            direction={direction || "row"}
            flexWrap="wrap"
            justifyContent="start"
            alignItems="center"
            gap={2}
        >
            {videos.map((item, idx) => (
                <Box key={idx}>
                    {item.id.videoId && <VideoCard video={item}/>}
                    {item.id.channelId && <ChannelCard channelDetail={item} />}
                </Box>
            ))}
        </Stack>
    )
}

export default Videos