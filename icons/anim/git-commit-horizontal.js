import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const before = el.querySelector('line[x1="3"]');
  animate(before, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const commit = el.querySelector("circle");
  animate(commit, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });

  const after = el.querySelector('line[x1="15"]');
  animate(after, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
