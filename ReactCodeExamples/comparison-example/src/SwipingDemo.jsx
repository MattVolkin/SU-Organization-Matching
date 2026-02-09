import { useState } from 'react'
import { useSprings, animated, to as interpolate } from '@react-spring/web'
import { useDrag } from '@use-gesture/react'

import "./SwipingDemo.css"

// using this https://use-gesture.netlify.app/docs/examples/ made a swiping gesture


// create data that allows for to and from functions to measure speed (and whether a card can be thrown or returned to the pile)
const to = (i) => ({
    x: 0, // because we only have left-right motion
    y: i * -4,
    scale: 1,
    rot: -10 + Math.random() *20, // add a bit of rotation so the cardsData are not perfectly flat
    delay: i * 100,
})

const from = (_i) => ({
    x:0, 
    rot:0, 
    scale: 1.5, 
    y: -1000 // far offscreen
})

const trans = (r, s) =>
  `perspective(1500px) rotateX(30deg) rotateY(${r / 10}deg) rotateZ(${r}deg) scale(${s})`


function Cards() {
    const [gone] = useState(() => new Set()) //set of cards that are passed
    const [props, api] = useSprings(cardsData.length, i => ({
        ...to(i),
        from: from(i),
    })) // create the cardsData and give them physics (using the springs we imported)

    //create a gesture to measure if its on top, the direction it is being tossed (important to determine a agree or disagree), distanced dragged, and velocity
    const gesture = useDrag(({ args: [index], active, movement: [mx], direction: [xDir], velocity: [vx] }) => {
        const trigger = vx > 0.2 // threshold for what is considered a swipe acceptable to clear card or return back to center
        if (!active && trigger) gone.add(index) // if the trigger velocity is reached (trigger) and it is not being currently held/dragged by the mouse (active), then add it to the gone list and remove it
        api.start(i => {
            if (index !== i) return //if we are not currently focused on the top card, don't worry about changing any spring data
            const isGone = gone.has(index) // check if card is in the gone list
            const x = isGone ? (200 + window.innerWidth) *xDir : active ? mx: 0 // when card is gone and not held, fly the card off screen, if the card is still being held, move back towards the center
            const rot = mx / 100 + (isGone ? xDir * 10 * vx : 0) // tilt the card depending on how fast the flick is
            const scale = active ? 1.1 : 1 // if a card is being held (active) then lift it up/make it bigger
            return {
                x,
                rot,
                scale,
                delay: undefined,
                config: { friction: 50, tension: active ? 800 : isGone ? 200 : 500 },
            }
        })
        if (!active && gone.size === cardsData.length) // if there no cards held and all cards are swiped, reset the pile. This would not be a functionality in the real product but for this example, reseting the list is ok
            setTimeout(() => {
                gone.clear() //remove cards from the removed/gone list
                api.start(i => to(i)) // animate them moving back to the center
            }, 600)
    })

    // now show the animation and cards

    return (
        <>
            {props.map(({x,y,rot,scale}, i) => (
                <animated.div className={cardsStyle} key={i} style={{ x, y }}>
                    {/*describe/form the visual of the card*/}
                    <animated.div 
                    {...gesture(i)}
                    style = {{
                        transform: interpolate([rot, scale], trans),
                        backgroundImage: 'cardsData[i].img'
                    }} 

                    />   
                </animated.div>

            ))} 

            <p>
                hello
            </p>
        
        
        </>
    )


    
}


function SwipingDemo() {
    return (
        <div className='Deck-Of-Cards'>
            <h1>this works</h1>
        </div>

    )

}


const cardsData = [
    {
        word: "Adaptable",
        desc: "Easily adjust to changes",
        img: '/images/lion.jpg',
    },
    
    {
        word: "Affectionate",
        desc: "Shows Love and care",
        img: '/images/raven.jpg',

    }

]

export default SwipingDemo;


