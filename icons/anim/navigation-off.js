import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [-2, 1, 0], translateY: [2, -1, 0], rotate: [-10, 4, 0], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector("line");
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
