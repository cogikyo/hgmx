import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const stick = el.querySelector("g");
  animate(stick, { rotate: [0, 8, -4, 0] }, { duration: τmid });
  const puck = el.querySelector("circle");
  animate(puck, { translateX: [0, 2.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
