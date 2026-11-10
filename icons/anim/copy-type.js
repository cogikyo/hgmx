import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 16"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const letter = el.querySelectorAll('path[d^="M12 13"], path[d="M15 12v6"]');
  animate(letter, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const foot = el.querySelector('path[d="M14 18h2"]');
  animate(foot, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
