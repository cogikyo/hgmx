import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const tickets = el.querySelectorAll('rect, path[d^="m3.173"]');
  animate(tickets, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const holes = ["M6 10V8", "M6 14v1", "M6 19v2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(holes, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
