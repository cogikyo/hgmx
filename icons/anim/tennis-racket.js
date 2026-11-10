import { animate, τmax } from "../motion.js";

export default function (el) {
  const strike = { duration: τmax, times: [0, 0.18, 0.36, 1] };
  const racket = el.querySelector("g");
  animate(racket, { rotate: [0, -8, 3, 0] }, strike);
  const ball = el.querySelector("circle");
  animate(ball, { translateX: [0, -1, -2, 0], translateY: [0, -1.25, -2.5, 0] }, strike);
}
