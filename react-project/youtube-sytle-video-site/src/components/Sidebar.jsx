import { Stack } from "@mui/system";
import { categories } from "../utils/constants";

const Sidebar = ({selectedCatagory, setSelectedCatagory}) => (
  <Stack 
    direction="row"
    sx={{
        overflowY: 'auto',
        height: {sx: 'auto', md: '95%'},
        flexDirection: { md: 'column'},
    }}
  >
      {categories.map((category) => (
          <button 
            className="category-btn"
            onClick={() => {
                setSelectedCatagory(category.name)
            }}
            style={{
                background: category.name === selectedCatagory && '$FC1503',
                color: 'white'
            }}
            key={category.name}
          >
              <span style={{
                  color: category.name === selectedCatagory ? 'white' : 'red',
                  marginRight: '15px',
              }}>{category.icon}</span>
              <span style={{
                  opacity: category.name === selectedCatagory ? '1' : '0.75',
              }}>{category.name}</span>
          </button>
      ))}
  </Stack>
)

export default Sidebar