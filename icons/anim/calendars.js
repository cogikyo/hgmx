import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("rect");
  animate(front, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const back = el.querySelectorAll('path[d^="M15.726"], path[d="M2 13h2"]');
  animate(back, { translateX: [4, -1, 0], translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
