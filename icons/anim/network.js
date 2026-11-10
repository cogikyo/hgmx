import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const links = el.querySelectorAll("path");
  animate(links, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const nodes = [el.querySelector('rect[x="2"]'), el.querySelector('rect[x="9"]'), el.querySelector('rect[x="16"]')];
  animate(nodes, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
