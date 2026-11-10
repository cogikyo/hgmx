import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const edges = el.querySelectorAll("path");
  animate(edges, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const at = (x, y) => el.querySelector(`circle[cx="${x}"][cy="${y}"]`);
  const nodes = [at(11, 4), at(20, 8), at(15, 20), at(4, 13)];
  animate(nodes, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
