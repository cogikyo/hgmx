import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("rect");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const keys = ["M8 10", "M12 10", "M16 10", "M8 14", "M12 14", "M8 18", "M12 18"].map((d) => el.querySelector(`path[d^="${d}"]`));
  keys.push(el.querySelector('line[y1="14"]'));
  animate(keys, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const display = el.querySelector('line[y1="6"]');
  animate(display, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
