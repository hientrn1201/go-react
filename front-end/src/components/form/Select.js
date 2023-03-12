const Select = (props) => {
  return (
    <div className="mb-3">
      <label htmlFor={props.className} className="form-label">
        {props.title}
      </label>
      <select 
        className="form-select"
        name={props.name}
        id={props.name}
        value={props.value}
        onChange={props.onChange}
      >
        <option value="">{props.placeHolder}</option>
        {props.options.map((option) => {
          return (
            <option key={option.id} value={option.id}>
              {option.value}
            </option>
          )
        })}
        <div className={props.errorDiv}>{props.errorMsg}</div>
      </select>
    </div>
  )
}

export default Select;