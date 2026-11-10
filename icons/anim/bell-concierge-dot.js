import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const button = el.querySelector("g g");
  animate(button, { scaleY: [1, 0.88, 1] }, { duration: τmin });
  const bell = el.querySelector("g");
  animate(bell, { rotate: [0, -6, 4, 0] }, { duration: τmid, delay: τmin });
  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
