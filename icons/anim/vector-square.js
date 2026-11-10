import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const at = (x, y) => el.querySelector(`circle[cx="${x}"][cy="${y}"]`);
  const nodes = [at(19, 5), at(5, 5), at(5, 19), at(19, 19)];
  animate(nodes, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const edges = el.querySelectorAll("path");
  animate(edges, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
