import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 16"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const slash = el.querySelector('path[d^="m18 12"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const dot = el.querySelector('path[d^="M12 18"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
