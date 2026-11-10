import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const forklift = el.querySelectorAll("path, circle");
  animate(forklift, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const forks = el.querySelector('path[d^="M15 19"]');
  animate(forks, { translateY: [0, -3, 0] }, { duration: τmid, delay: τmin });
}
