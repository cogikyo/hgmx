import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const dishwasher = el.querySelector("rect");
  animate(dishwasher, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const streaks = el.querySelectorAll('path[d="m9 11-2 7"], path[d="m13 11-2 7"]');
  animate(streaks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
