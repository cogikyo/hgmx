import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const edges = el.querySelectorAll("path");
  animate(edges, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const commits = [el.querySelector('circle[cy="6"][cx="5"]'), el.querySelector('circle[cy="18"]'), el.querySelector('circle[cx="19"]')];
  animate(commits, { scale: [1, 1.25, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
