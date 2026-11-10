import { animate, stagger } from "../motion.js";

export default function (el) {
  const dashes = ["M12 5V3", "M12 9v3", "M12 17v4"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(dashes, { scaleY: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
