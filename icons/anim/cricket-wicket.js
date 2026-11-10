import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ball = el.querySelector("circle");
  animate(ball, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid });
  const bails = el.querySelectorAll('path[d="m6 2 4 2"], path[d="m14 3 4-1"]');
  animate(bails, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
