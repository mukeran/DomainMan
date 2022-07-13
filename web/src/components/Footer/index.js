import React from 'react'
import './style.less'

class Footer extends React.Component {
  render () {
    return (
      <span>
        <b>DomainMan Web UI</b>, Part of <a href="https://github.com/mukeran/DomainMan">DomainMan</a> Project<br />
        Copyright &copy; 2020-2022 <a href="https://blog.mkr.im/devops">mukeran</a>
      </span>
    )
  }
}

export default Footer