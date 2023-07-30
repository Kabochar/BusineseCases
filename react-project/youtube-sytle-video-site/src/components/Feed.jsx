import { Box, Stack, Typography } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { fetchFromAPI } from '../utils/fetchFromAPI';
import { Sidebar, Videos } from './';

const Feed = () => {
	const [selectedCatagory, setSelectedCatagory] = useState('New')
	const [videos, setVideos] = useState([]);

	useEffect(() => {
		setVideos([]);

		fetchFromAPI(`search?part=snippet&q=${selectedCatagory}`).then((data) => {
			setVideos(data.items)
		});
	}, [selectedCatagory])
    
	return (
		<Stack sx={{
			flexDirection: {
				sx: 'column',
				md: 'row'
			}
		}}>
			<Box sx={{
				height: {
					sx: 'auto',
					md: '92vh'
				},
				borderRight: '1px solid #3d3d3d',
				px: { sx: 0, md: 2 }
			}}>
				<Sidebar
					selectedCatagory={selectedCatagory}
					setSelectedCatagory={setSelectedCatagory}
				/>
				<Typography
					className="copyright"
					variant="body2"
					sx={{
						mt: 1.5,
						color: '#fff'
					}}>
					Copyright © 2023 kabochar
				</Typography>
			</Box>
			<Box p={2} sx={{ overflowY: 'auto', height: '90vh', flex: 2 }}>
				<Typography variant="h4"
					fontWeight="bold"
					mb={2}
					sx={{
						color: 'white'
					}}
				>
					{selectedCatagory} <span style={{ color: '#F31503' }}>videos</span>
				</Typography>

				<Videos videos={videos} />
			</Box>
		</Stack>
	)
};

export default Feed;
