import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const paddle = el.querySelector("g");
  animate(paddle, { rotate: [0, -10, 4, 0] }, { duration: τmid });

  const ball = el.querySelector("circle");
  animate(ball, { translateX: [0, 2.5, -0.5, 0], translateY: [0, 2.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
