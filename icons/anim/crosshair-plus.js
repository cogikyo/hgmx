import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const ticks = el.querySelector("g");
  animate(ticks, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
  animate(ticks, { scale: [1, 0.75, 1] }, { duration: τmid, delay: τmin });
}
