import React, { useState, useEffect, useMemo } from 'react';
import {
  Container,
  Grid,
  Card,
  CardContent,
  CardMedia,
  Typography,
  Button,
  Box,
  CircularProgress,
  AppBar,
  Toolbar,
  TextField,
  Chip,
  Stack,
  IconButton,
  Fade,
  Paper,
  Divider,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import ViewModuleIcon from '@mui/icons-material/ViewModule';
import ViewListIcon from '@mui/icons-material/ViewList';
import SearchIcon from '@mui/icons-material/Search';
import FilterListIcon from '@mui/icons-material/FilterList';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import axios from 'axios';

const getTheme = (mode) => createTheme({
  palette: {
    mode,
    ...(mode === 'dark' ? {
      primary: {
        main: '#90caf9',
      },
      secondary: {
        main: '#f48fb1',
      },
      background: {
        default: '#121212',
        paper: '#1e1e1e',
      },
    } : {
      primary: {
        main: '#1976d2',
      },
      secondary: {
        main: '#dc004e',
      },
    }),
  },
  components: {
    MuiCard: {
      styleOverrides: {
        root: {
          transition: 'transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
          '&:hover': {
            transform: 'translateY(-4px)',
            boxShadow: mode === 'dark' 
              ? '0 8px 16px rgba(0,0,0,0.4)'
              : '0 8px 16px rgba(0,0,0,0.1)',
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          transition: 'all 0.2s ease-in-out',
          '&:hover': {
            transform: 'scale(1.05)',
          },
        },
      },
    },
  },
});

function App() {
  const [videos, setVideos] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [cursor, setCursor] = useState('');
  const [hasMore, setHasMore] = useState(true);
  const [startDate, setStartDate] = useState(null);
  const [endDate, setEndDate] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedChannels, setSelectedChannels] = useState([]);
  const [mode, setMode] = useState(() => {
    const savedMode = localStorage.getItem('themeMode');
    return savedMode || 'dark';
  });
  const [viewMode, setViewMode] = useState(() => {
    const savedViewMode = localStorage.getItem('viewMode');
    return savedViewMode || 'grid';
  });

  useEffect(() => {
    localStorage.setItem('themeMode', mode);
  }, [mode]);

  useEffect(() => {
    localStorage.setItem('viewMode', viewMode);
  }, [viewMode]);

  const toggleColorMode = () => {
    setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'));
  };

  const toggleViewMode = () => {
    setViewMode((prevMode) => (prevMode === 'grid' ? 'list' : 'grid'));
  };

  const uniqueChannels = useMemo(() => {
    const channels = new Set(videos.map(video => video.channel_title));
    return Array.from(channels).sort();
  }, [videos]);

  const fetchVideos = async (nextCursor = '') => {
    try {
      setLoading(true);
      setError(null);
      
      const params = new URLSearchParams({
        limit: 12,
      });
      
      if (nextCursor) {
        params.append('cursor', nextCursor);
      }
      
      const response = await axios.get(`http://localhost:8080/api/v1/videos?${params}`);
      const { videos: newVideos, next_cursor, has_more } = response.data;
      
      if (nextCursor) {
        setVideos(prev => [...prev, ...newVideos]);
      } else {
        setVideos(newVideos);
      }
      
      setCursor(next_cursor);
      setHasMore(has_more);
    } catch (err) {
      setError('Failed to fetch videos. Please try again.');
      console.error('Error fetching videos:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchVideos();
  }, []);

  const handleLoadMore = () => {
    if (cursor && hasMore) {
      fetchVideos(cursor);
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const handleChannelToggle = (channel) => {
    setSelectedChannels(prev => 
      prev.includes(channel)
        ? prev.filter(c => c !== channel)
        : [...prev, channel]
    );
  };

  const filteredVideos = useMemo(() => {
    return videos.filter(video => {
      const matchesSearch = searchTerm === '' || 
        video.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        video.description.toLowerCase().includes(searchTerm.toLowerCase());

      const videoDate = new Date(video.published_at);
      const matchesStartDate = !startDate || videoDate >= startDate;
      const matchesEndDate = !endDate || videoDate <= endDate;

      const matchesChannel = selectedChannels.length === 0 || 
        selectedChannels.includes(video.channel_title);

      return matchesSearch && matchesStartDate && matchesEndDate && matchesChannel;
    });
  }, [videos, searchTerm, startDate, endDate, selectedChannels]);

  const renderVideoCard = (video) => {
    if (viewMode === 'list') {
      return (
        <Fade in timeout={500}>
          <Card sx={{ 
            display: 'flex', 
            mb: 2,
            overflow: 'hidden',
            borderRadius: 2,
          }}>
            <CardMedia
              component="img"
              sx={{ 
                width: 320, 
                height: 180,
                objectFit: 'cover',
              }}
              image={video.thumbnail_url}
              alt={video.title}
            />
            <CardContent sx={{ 
              flex: 1,
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'space-between',
            }}>
              <Box>
                <Typography gutterBottom variant="h6" component="h2" sx={{ 
                  fontWeight: 600,
                  lineHeight: 1.3,
                }}>
                  {video.title}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ 
                  mb: 2,
                  display: '-webkit-box',
                  WebkitLineClamp: 3,
                  WebkitBoxOrient: 'vertical',
                  overflow: 'hidden',
                }}>
                  {video.description}
                </Typography>
              </Box>
              <Box>
                <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                  Published: {formatDate(video.published_at)}
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                  Channel: {video.channel_title}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Fade>
      );
    }

    return (
      <Fade in timeout={500}>
        <Card sx={{ 
          height: '100%', 
          display: 'flex', 
          flexDirection: 'column',
          borderRadius: 2,
        }}>
          <CardMedia
            component="img"
            height="200"
            image={video.thumbnail_url}
            alt={video.title}
            sx={{ objectFit: 'cover' }}
          />
          <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
            <Typography gutterBottom variant="h6" component="h2" sx={{ 
              fontWeight: 600,
              lineHeight: 1.3,
            }}>
              {video.title}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ 
              mb: 2,
              display: '-webkit-box',
              WebkitLineClamp: 3,
              WebkitBoxOrient: 'vertical',
              overflow: 'hidden',
            }}>
              {video.description}
            </Typography>
            <Box sx={{ mt: 'auto' }}>
              <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                Published: {formatDate(video.published_at)}
              </Typography>
              <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                Channel: {video.channel_title}
              </Typography>
            </Box>
          </CardContent>
        </Card>
      </Fade>
    );
  };

  return (
    <ThemeProvider theme={getTheme(mode)}>
      <CssBaseline />
      <LocalizationProvider dateAdapter={AdapterDayjs}>
        <Box sx={{ flexGrow: 1 }}>
          <AppBar position="static" elevation={0} sx={{ 
            background: mode === 'dark' 
              ? 'linear-gradient(45deg, #1a237e 30%, #0d47a1 90%)'
              : 'linear-gradient(45deg, #1976d2 30%, #2196f3 90%)',
          }}>
            <Toolbar>
              <Typography variant="h6" component="div" sx={{ 
                flexGrow: 1,
                fontWeight: 600,
                letterSpacing: '0.5px',
              }}>
                YouTube Videos Dashboard
              </Typography>
              <IconButton onClick={toggleViewMode} color="inherit" sx={{ mr: 1 }}>
                {viewMode === 'grid' ? <ViewListIcon /> : <ViewModuleIcon />}
              </IconButton>
              <IconButton onClick={toggleColorMode} color="inherit">
                {mode === 'dark' ? <Brightness7Icon /> : <Brightness4Icon />}
              </IconButton>
            </Toolbar>
          </AppBar>

          <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Paper elevation={0} sx={{ 
              p: 3, 
              mb: 4,
              borderRadius: 2,
              background: mode === 'dark' ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.02)',
            }}>
              <Grid container spacing={3}>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    label="Search Videos"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="Search in title or description..."
                    InputProps={{
                      startAdornment: <SearchIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                    }}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <DatePicker
                    label="Start Date"
                    value={startDate}
                    onChange={setStartDate}
                    slotProps={{ textField: { fullWidth: true } }}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <DatePicker
                    label="End Date"
                    value={endDate}
                    onChange={setEndDate}
                    slotProps={{ textField: { fullWidth: true } }}
                  />
                </Grid>
              </Grid>

              <Box sx={{ mt: 3 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  <FilterListIcon sx={{ mr: 1, color: 'text.secondary' }} />
                  <Typography variant="subtitle1" sx={{ fontWeight: 500 }}>
                    Filter by Channel
                  </Typography>
                </Box>
                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                  {uniqueChannels.map((channel) => (
                    <Chip
                      key={channel}
                      label={channel}
                      onClick={() => handleChannelToggle(channel)}
                      color={selectedChannels.includes(channel) ? "primary" : "default"}
                      sx={{ m: 0.5 }}
                    />
                  ))}
                </Stack>
              </Box>

              <Box sx={{ mt: 3, display: 'flex', gap: 1 }}>
                <Button
                  variant="contained"
                  onClick={() => {
                    setSearchTerm('');
                    setStartDate(null);
                    setEndDate(null);
                    setSelectedChannels([]);
                  }}
                  sx={{
                    background: mode === 'dark' 
                      ? 'linear-gradient(45deg, #f48fb1 30%, #f06292 90%)'
                      : 'linear-gradient(45deg, #dc004e 30%, #f50057 90%)',
                    '&:hover': {
                      background: mode === 'dark'
                        ? 'linear-gradient(45deg, #f06292 30%, #ec407a 90%)'
                        : 'linear-gradient(45deg, #f50057 30%, #ff4081 90%)',
                    },
                  }}
                >
                  Clear Filters
                </Button>
              </Box>
            </Paper>

            <Grid container spacing={3}>
              {error && (
                <Grid item xs={12}>
                  <Typography color="error">{error}</Typography>
                </Grid>
              )}

              {viewMode === 'grid' ? (
                filteredVideos.map((video) => (
                  <Grid item xs={12} sm={6} md={4} key={video.id}>
                    {renderVideoCard(video)}
                  </Grid>
                ))
              ) : (
                <Grid item xs={12}>
                  {filteredVideos.map((video) => (
                    <Box key={video.id}>
                      {renderVideoCard(video)}
                    </Box>
                  ))}
                </Grid>
              )}

              {loading && (
                <Grid item xs={12} sx={{ textAlign: 'center', my: 4 }}>
                  <CircularProgress />
                </Grid>
              )}

              {hasMore && !loading && (
                <Grid item xs={12} sx={{ textAlign: 'center', my: 4 }}>
                  <Button
                    variant="contained"
                    onClick={handleLoadMore}
                    disabled={loading}
                    sx={{
                      px: 4,
                      py: 1.5,
                      background: mode === 'dark'
                        ? 'linear-gradient(45deg, #90caf9 30%, #64b5f6 90%)'
                        : 'linear-gradient(45deg, #1976d2 30%, #2196f3 90%)',
                      '&:hover': {
                        background: mode === 'dark'
                          ? 'linear-gradient(45deg, #64b5f6 30%, #42a5f5 90%)'
                          : 'linear-gradient(45deg, #2196f3 30%, #1e88e5 90%)',
                      },
                    }}
                  >
                    Load More
                  </Button>
                </Grid>
              )}
            </Grid>
          </Container>
        </Box>
      </LocalizationProvider>
    </ThemeProvider>
  );
}

export default App; 