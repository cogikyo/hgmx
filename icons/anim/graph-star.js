import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const spokes = el.querySelectorAll("path");
  animate(spokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const nodes = ['circle[cy="4"][cx="15"]', 'circle[cx="20"]', 'circle[cx="16"]', 'circle[cy="19"]', 'circle[cy="4"][cx="4"]'].map((s) => el.querySelector(s));
  animate(nodes, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
