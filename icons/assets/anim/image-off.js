import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelectorAll('path[d^="M3.59"], path[d^="M21 15"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('line[x1="2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector('path[d^="M10.41"]');
  animate(sun, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
