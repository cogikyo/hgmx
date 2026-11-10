import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const halves = [el.querySelector('path[d^="M3.712"]'), el.querySelector('path[d^="M20.288"]')];
  animate(halves, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dashes = ["M12 2v2", "M12 8v2", "M12 14v2", "M12 20v2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(dashes, { scaleY: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
