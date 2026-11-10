import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M20 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const before = el.querySelector('path[d="M7 13h3"]');
  animate(before, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const commit = el.querySelector("circle");
  animate(commit, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });

  const after = el.querySelector('path[d="M14 13h3"]');
  animate(after, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
